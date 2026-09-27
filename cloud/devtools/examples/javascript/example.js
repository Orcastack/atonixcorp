import compute from "../../sdk/javascript/compute.js";
import identity from "../../sdk/javascript/identity.js";
import telemetry from "../../sdk/javascript/telemetry.js";

process.env.ATCLOUD_ENDPOINT = "https://api.atcloud.africa/v1";
process.env.ATCLOUD_TOKEN = "YOUR_TOKEN";

async function main() {
  console.log("Compute:", await compute.list());
  console.log("Identity:", await identity.info());
  console.log("Telemetry:", await telemetry.snapshot());
}

main();
