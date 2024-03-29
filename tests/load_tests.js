import { expect } from "https://jslib.k6.io/k6chaijs/4.3.4.1/index.js";
import { randomIntBetween } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";
import http from "k6/http";
import { sleep } from "k6";

export const options = {
  stages: [
    // simulate ramp-up of traffic from 1 to 100 users over 5 seconds
    { duration: "5s", target: 100 },
    // stay at 100 users for 10 seconds
    { duration: "10s", target: 100 },
    // ramp-down to 0 users over 5 seconds
    { duration: "5s", target: 0 },
  ],
  thresholds: {
    // the rate of HTTP errors should be less than 1%
    http_req_failed: ["rate < 0.01"],
    // 90% of requests must finish within 400ms, 95% within 800, and 99.9% within 2s
    http_req_duration: ["p(90) < 400", "p(95) < 800", "p(99.9) < 2000"],
    // the rate of successful checks should be higher than 90%
    checks: ["rate > 0.9"],
  },
};

function expectBody(response) {
  return expect(response.json(), "response body");
}

function expectBodyProperty(response, property) {
  return expect(response.json(property), `property '${property}'`);
}

export default function () {
  const throws = randomIntBetween(1, 100);
  const faces = randomIntBetween(2, 100);

  const response = http.get(
    `http://${__ENV.SERVICE_ADDRESS}/api/v1/dice?throws=${throws}&faces=${faces}`
  );

  expect(response.status, "status code").to.equal(200);
  expect(response).to.have.validJsonBody();

  expectBody(response).to.have.property("Throws");
  expectBodyProperty(response, "Throws")
    .to.be.an("array")
    .to.have.lengthOf(throws);

  expectBody(response).to.have.property("Statistics");

  expectBodyProperty(response, "Statistics").to.have.property("Minimum");
  expectBodyProperty(response, "Statistics.Minimum").to.equal(
    Math.min(...response.json("Throws"))
  );

  expectBodyProperty(response, "Statistics").to.have.property("Maximum");
  expectBodyProperty(response, "Statistics.Maximum").to.equal(
    Math.max(...response.json("Throws"))
  );

  expectBodyProperty(response, "Statistics").to.have.property("Sum");
  expectBodyProperty(response, "Statistics.Sum").to.equal(
    response.json("Throws").reduce((result, item) => result + item, 0)
  );

  sleep(1);
}
