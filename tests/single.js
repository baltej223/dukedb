import Duke from "duke-client";

const db = new Duke();

db.seturl("duke://localhost:9000");

console.log("Connecting...");

const connected = await db.connect();

if (!connected) {
  console.error("Failed to connect.");
  process.exit(1);
}

console.log("Connected.");

await db.PUT("hello", "world");

const value = await db.GET("hello");

if (value !== "world") {
  console.error("GET returned wrong value:", value);
  process.exit(1);
}

console.log("Single node test passed.");
