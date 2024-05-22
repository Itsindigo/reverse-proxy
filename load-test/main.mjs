import fetch from "node-fetch";

// write a generic `request` function using node-fetch, return all responses, do not throw >= 400, assume application/json always
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

const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

async function main() {
  // fire requests constantly for 60 seconds, store bad requests in one array, good requests in another
  const badRequests = [];
  const goodRequests = [];
  const startTime = Date.now();
  const endTime = startTime + 60000;
  while (Date.now() < endTime) {
    // make GET request to http://localhost:6666/api/hello", if status code == 200, store in goodRequests, if 429 store in badRequests
    const response = await request("http://localhost:6666/api/hello", "GET");
    if (response.statusCode === 200) {
      goodRequests.push(response);
    } else if (response.statusCode === 429) {
      badRequests.push(response);
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
