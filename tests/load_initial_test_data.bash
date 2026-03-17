#!/usr/bin/env bash

set -euo pipefail

base_url="${TAIGO_BASE_URL:-http://localhost:9000}"
base_url="${base_url%/}"
primary_username="${TAIGO_USERNAME:-admin}"
primary_password="${TAIGO_PASSWORD:-123123}"

role_matrix_username="${TAIGO_MEMBER_USERNAME:-taigo_role_ci_r4h7m2}"
role_matrix_password="${TAIGO_MEMBER_PASSWORD:-N7k3Q9x2R5m8P1c4}"
role_matrix_email="${TAIGO_MEMBER_EMAIL:-${role_matrix_username}@example.invalid}"
role_matrix_full_name="${TAIGO_MEMBER_FULL_NAME:-Taigo Role Matrix}"

wait_for_taiga() {
	while :; do
		if curl --silent --show-error --fail --output /dev/null "${base_url}/api/v1/stats/discover"; then
			break
		fi
		sleep 1
	done
}

post_json() {
	local endpoint="$1"
	local payload="$2"
	local report_failure_body="${3:-1}"
	local body_file
	body_file="$(mktemp)"

	local status
	status="$(
		curl \
			--silent \
			--show-error \
			--output "${body_file}" \
			--write-out "%{http_code}" \
			--header "Content-Type: application/json" \
			--request POST \
			--data "${payload}" \
			"${base_url}${endpoint}"
	)"

	if [[ "${status}" != "200" && "${status}" != "201" && "${report_failure_body}" == "1" ]]; then
		cat "${body_file}" >&2
	fi
	rm -f "${body_file}"
	printf '%s' "${status}"
}

verify_primary_test_credentials() {
	local login_payload
	login_payload=$(
		cat <<EOF
{"type":"normal","username":"${primary_username}","password":"${primary_password}"}
EOF
	)

	local login_status
	login_status="$(post_json "/api/v1/auth" "${login_payload}")"
	echo
	if [[ "${login_status}" != "200" ]]; then
		echo "failed to authenticate primary test user ${primary_username}; HTTP ${login_status}" >&2
		return 1
	fi

	echo "Verified primary test credentials: ${primary_username}"
}

ensure_role_matrix_user() {
	local login_payload
	login_payload=$(
		cat <<EOF
{"type":"normal","username":"${role_matrix_username}","password":"${role_matrix_password}"}
EOF
	)

	local login_status
	login_status="$(post_json "/api/v1/auth" "${login_payload}" 0)"
	if [[ "${login_status}" == "200" ]]; then
		echo
		echo "Role-matrix user already exists: ${role_matrix_username}"
		return 0
	fi
	echo

	local register_payload
	register_payload=$(
		cat <<EOF
{"type":"public","username":"${role_matrix_username}","password":"${role_matrix_password}","email":"${role_matrix_email}","full_name":"${role_matrix_full_name}","accepted_terms":true}
EOF
	)

	local register_status
	register_status="$(post_json "/api/v1/auth/register" "${register_payload}")"
	echo
	case "${register_status}" in
		200|201)
			echo "Created role-matrix user: ${role_matrix_username}"
			;;
		*)
			echo "failed to create role-matrix user ${role_matrix_username}; HTTP ${register_status}" >&2
			return 1
			;;
	esac
}

wait_for_taiga
echo "Taiga is up at ${base_url}"

tgback="$(docker ps --filter "name=taiga-back" --filter "status=running" --format '{{.Names}}')"
if [[ -z "${tgback}" ]]; then
	echo "could not find a running taiga-back container" >&2
	exit 1
fi

docker cp initial_test_data.json "${tgback}:/taiga-back/media/initial_test_data.json"

cd taiga-docker
docker compose -f docker-compose.yml -f docker-compose-inits.yml run --rm taiga-manage loaddata /taiga-back/media/initial_test_data.json
cd ..

verify_primary_test_credentials
ensure_role_matrix_user

cat <<EOF
Seeded role-matrix credentials:
  username: ${role_matrix_username}
  password: ${role_matrix_password}
  expectation: forbid
EOF
