// scripts/load-test.k6.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 50 },
    { duration: '1m', target: 100 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<300'], // 95-й перцентиль < 300 мс
  },
};

const BASE_URL = 'http://localhost:8080';

export default function () {
  // создаём команду и пользователей один раз (в реальности можно в setup)
  let teamRes = http.post(`${BASE_URL}/teams`, JSON.stringify({name: "loadteam"}), {headers: {'Content-Type': 'application/json'}});
  let teamId = JSON.parse(teamRes.body).id;

  let aliceRes = http.post(`${BASE_URL}/users`, JSON.stringify({username: `alice${__VU}`, team_id: teamId}), {headers: {'Content-Type': 'application/json'}});
  let bobRes = http.post(`${BASE_URL}/users`, JSON.stringify({username: `bob${__VU}`, team_id: teamId}), {headers: {'Content-Type': 'application/json'}});

  let aliceId = JSON.parse(aliceRes.body).id;

  // создаём PR
  let prRes = http.post(`${BASE_URL}/pull-requests`, JSON.stringify({title: "load test", author_id: aliceId}), {headers: {'Content-Type': 'application/json'}});
  let prId = JSON.parse(prRes.body).id;

  check(prRes, { 'PR created': (r) => r.status === 200 });

  sleep(0.3);
}