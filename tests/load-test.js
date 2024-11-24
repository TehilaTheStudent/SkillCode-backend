import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 10, // Number of virtual users
  duration: '30s', // Duration of the test
};

export default function () {
  const url = 'http://localhost:8080/skillcode/questions/674278c2229903f6618ec65e/test';

  // JSON payload
  const payload = JSON.stringify({
    language: "Python",
    code: "binary_search",
  });

  // Request headers
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  // Send POST request
  const response = http.post(url, payload, params);

  // Validate response
  check(response, {
    'status is 200': (r) => r.status === 200,
    'response time < 200ms': (r) => r.timings.duration < 200,
  });

  // Simulate a pause between requests
  sleep(1);
}
