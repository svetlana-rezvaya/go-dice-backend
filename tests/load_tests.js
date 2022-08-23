import http from "k6/http";
import { sleep } from "k6";

export const options = {
  stages: [
    { duration: "5s", target: 100 }, // simulate ramp-up of traffic from 1 to 100 users over 5 seconds
    { duration: "10s", target: 100 }, // stay at 100 users for 10 seconds
    { duration: "5s", target: 0 }, // ramp-down to 0 users over 5 seconds
  ],
  thresholds: {
    http_req_failed: ["rate < 0.01"], // HTTP errors should be less than 1%
    // 90% of requests must finish within 400ms, 95% within 800, and 99.9% within 2s.
    http_req_duration: ["p(90) < 400", "p(95) < 800", "p(99.9) < 2000"],
  },
};

export default function () {
  http.get("https://test.k6.io");
  sleep(1);
}
