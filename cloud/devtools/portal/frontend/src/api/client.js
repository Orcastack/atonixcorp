const API = import.meta.env.VITE_ATCLOUD_API;

export async function apiGet(path) {
  const res = await fetch(API + path, {
    headers: {
      Authorization: "Bearer " + localStorage.getItem("token"),
    },
  });

  if (!res.ok) throw new Error(await res.text());
  return res.json();
}
