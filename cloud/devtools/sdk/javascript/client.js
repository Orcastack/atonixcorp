const API_ENDPOINT = process.env.ATCLOUD_ENDPOINT;
const API_TOKEN = process.env.ATCLOUD_TOKEN;

async function request(method, path, body) {
  const url = API_ENDPOINT + path;

  const options = {
    method,
    headers: {
      "Authorization": `Bearer ${API_TOKEN}`,
      "Content-Type": "application/json"
    }
  };

  if (body) {
    options.body = JSON.stringify(body);
  }

  const res = await fetch(url, options);

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`API Error ${res.status}: ${text}`);
  }

  return res.json();
}

export default {
  request
};
