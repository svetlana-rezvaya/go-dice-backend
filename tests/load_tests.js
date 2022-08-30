import http from "k6/http";
import { check, sleep } from "k6";

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
    // the rate of successful checks should be higher than 90%
    checks: ["rate > 0.9"],
  },
};

function getRandomNumberInRange(minimum, maximum) {
  const randomNumber = Math.random(); // return a random number in the range 0 to less than 1 (inclusive of 0, but not 1)
  // transform range [0, 1) to range [minimum, maximum)
  const transformedRandomNumber = (maximum - minimum) * randomNumber + minimum;
  return Math.round(transformedRandomNumber);
}

export default function () {
  const throws = getRandomNumberInRange(1, 100);
  const faces = getRandomNumberInRange(2, 100);
  const response = http.get(
    `http://go-dice-backend:9090/api/v1/dice?throws=${throws}&faces=${faces}`
  );
  check(response, {
    "is status 200": (response) => response.status === 200,
    "has Throws array": (response) => {
      // try to parse the response as JSON and receive the Throws field from it
      const receivedThrows = response.json("Throws");
      // check if received throws is an array and its length equals the desired length
      return Array.isArray(receivedThrows) && receivedThrows.length === throws;
    },
  });

  sleep(1);
}
