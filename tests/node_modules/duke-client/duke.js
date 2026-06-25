function dukeToHttp(url) {
  return url.replace("duke://", "http://");
}

export default class Duke {
  constructor() {
    this.baseObject = {
      baseURLs: [],
      workingURLs: [],
    };
  }

  seturl(...urls) {
    urls.forEach((url) => {
      if (!url.startsWith("duke://")) {
        throw new Error(`Invalid URL: ${url}`);
      }
    });

    this.baseObject.baseURLs.push(...urls);
  }

  check() {
    if (this.baseObject.baseURLs.length === 0) {
      throw new Error("No URLs defined.");
    }

    if (this.baseObject.workingURLs.length === 0) {
      throw new Error(
        "No healthy Duke nodes available. Call connect() first.",
      );
    }
  }

  chooseRandomWorkingUrl() {
    const len = this.baseObject.workingURLs.length;

    if (len === 0) {
      throw new Error("No working URLs available.");
    }

    const randomIndex = Math.floor(Math.random() * len);

    const urlIndex = this.baseObject.workingURLs[randomIndex];

    return this.baseObject.baseURLs[urlIndex];
  }

  async check_health(urlIndex) {
    if (
      urlIndex < 0 ||
      urlIndex >= this.baseObject.baseURLs.length
    ) {
      throw new Error(
        "Invalid URL index provided.",
      );
    }

    try {
      const response = await fetch(
        dukeToHttp(this.baseObject.baseURLs[urlIndex]) + "/health",
      );

      const data = await response.text();

      return data === "OK";
    } catch {
      return false;
    }
  }

  async connect() {
    console.log("connecting");
    this.baseObject.workingURLs = [];

    const total =
      this.baseObject.baseURLs.length;

    let healthy = 0;

    for (let i = 0; i < total; i++) {
      const ok = await this.check_health(i);

      if (ok) {
        this.baseObject.workingURLs.push(i);
        healthy++;
      }
    }
    console.log("Connected!")
    return healthy > 0;

  }

  async GET(key) {
    this.check();

    const url = dukeToHttp(
      this.chooseRandomWorkingUrl(),
    );

    const response = await fetch(`${url}/get?key=${encodeURIComponent(key)}`);

    const data = await response.json();

    if (data.found) {
      return data.value;
    }

    throw new Error(data.error || "Key not found.");
  }

  async PUT(key, value) {
    if (value == undefined || key == undefined) {
      throw new Error(key == undefined ? "Key Value can't be undefined!" : value == undefined ? "Value can't be undefined for PUT." : "");
    }
    this.check();

    const url = dukeToHttp(this.chooseRandomWorkingUrl());

    const response = await fetch(`${url}/put`,
      {
        method: "POST",
        headers: {
          "Content-Type":
            "application/json",
        },
        body: JSON.stringify({
          key,
          value,
        }),
      },
    );

    const data = await response.json();

    if (!data.success) {
      throw new Error(data.error || "PUT failed.");
    }

    return true;
  }

  async batch_GET(keysArr, BATCH_SIZE) {
    let response = [];
    for (let i = 0; i < keysArr.length; i += BATCH_SIZE) {
      let batch = [];

      for (let j = i; j < i + BATCH_SIZE; j++) {
        batch.push(
          db.GET(keysArr[j])
        );
      }
      let batch_response = await Promise.all(batch);
      response = [...response, ...batch_response]
    }
    return response
  }


  async batch_PUT(keysArr, valuesArr, BATCH_SIZE) {
    if (keysArr.length != valuesArr.length) {
      throw new Error("Different length of keys array and value arrays provided.");
    }
    let response = [];
    for (let i = 0; i < keysArr.length; i += BATCH_SIZE) {
      let batch = [];

      for (let j = i; j < i + BATCH_SIZE; j++) {
        batch.push(
          db.PUT(keysArr[j], valuesArr[j])
        );
      }
      let batch_response = await Promise.all(batch);
      response = [...response, ...batch_response]
    }
    return response
  }
}
