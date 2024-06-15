import { check } from "k6";
import http from "k6/http";

export default function () {
  let res = http.get("http://gw:8080/api/rate");
  check(res, {
    "is status 200": (r) => r.status === 200,
  });
}
