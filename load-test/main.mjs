import fetch from "node-fetch";

async function request(url, method, body) {
  const response = await fetch(url, {
    method,
    body: JSON.stringify(body),
    headers: {
      "Content-Type": "application/json",
    },
  });

  return {
    statusCode: response.status,
  };
}

async function main() {
  const badRequests = [];
  const goodRequests = [];
  const startTime = Date.now();
  const endTime = startTime + 60000;
  while (Date.now() < endTime) {
    const helloResponse = await request(
      "http://localhost:6666/api/hello",
      "GET"
    );
    if (helloResponse.statusCode === 200) {
      goodRequests.push(helloResponse);
    } else if (helloResponse.statusCode === 429) {
      badRequests.push(helloResponse);
    }

    const goodbyeResponse = await request(
      "http://localhost:6666/api/goodbye",
      "GET"
    );
    if (goodbyeResponse.statusCode === 200) {
      goodRequests.push(goodbyeResponse);
    } else if (goodbyeResponse.statusCode === 429) {
      badRequests.push(goodbyeResponse);
    }
  }

  return {
    successCount: goodRequests.length,
    limitedCount: badRequests.length,
  };
}

main()
  .then((results) => {
    console.log({ results });
  })
  .catch((err) => {
    console.error(err);
  });
