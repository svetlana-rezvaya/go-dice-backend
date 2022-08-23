import http from "k6/http";
import { sleep } from "k6";

export const options = {
  stages: [
    { duration: "5s", target: 100 }, // simulate ramp-up of traffic from 1 to 100 users over 5 seconds
    { duration: "10s", target: 100 }, // stay at 100 users for 10 seconds
    { duration: "5s", target: 0 }, // ramp-down to 0 users over 5 seconds
  ],
};

export default function () {
  http.get("https://test.k6.io");
  sleep(1);
}
