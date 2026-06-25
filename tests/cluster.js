import Duke from "duke-client";

const db = new Duke();

db.seturl(
  "duke://localhost:9000",
  "duke://localhost:9001",
  "duke://localhost:9002",
  "duke://localhost:9003",
  "duke://localhost:9004",
);

console.log("Connecting...");

const connected = await db.connect();

if (!connected) {
  console.error("Failed to connect.");
  process.exit(1);
}

console.log("Connected.");

// ---------------------
// PUT TEST
// ---------------------

const COUNT = 1000;

for (let i = 0; i < COUNT; i++) {
  try {
    await db.PUT(
      `key-${i}`,
      `value-${i}`,
    );
  } catch (err) {
    console.error("PUT failed:", err);
    process.exit(1);
  }
}

console.log(`Inserted ${COUNT} keys.`);

// ---------------------
// GET TEST
// ---------------------

for (let i = 0; i < COUNT; i++) {
  try {
    const value = await db.GET(`key-${i}`);

    if (value !== `value-${i}`) {
      console.error(
        `Incorrect value for key-${i}. Expected value-${i}, got ${value}`,
      );
      process.exit(1);
    }
  } catch (err) {
    console.error(`GET failed for key-${i}:`, err);
    process.exit(1);
  }
}

console.log(`Verified ${COUNT} keys.`);

console.log("✅ Cluster integration test passed.");
