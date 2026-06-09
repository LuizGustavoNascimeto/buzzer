const apiBaseUrl = process.env.REACT_APP_BACKEND_URL || "http://localhost:8080";

export async function fetchActivities() {
  const res = await fetch(`${apiBaseUrl}/api/activities/home`);
  return res.json();
}

function ttlToExpiresAt(ttl) {
  const now = new Date();

  const map = {
    "30-days": 30 * 24 * 60 * 60 * 1000,
    "7-days": 7 * 24 * 60 * 60 * 1000,
    "3-days": 3 * 24 * 60 * 60 * 1000,
    "1-day": 1 * 24 * 60 * 60 * 1000,
    "12-hours": 12 * 60 * 60 * 1000,
    "3-hours": 3 * 60 * 60 * 1000,
    "1-hour": 1 * 60 * 60 * 1000,
  };

  const ms = map[ttl];

  return ms ? new Date(now.getTime() + ms).toISOString() : null;
}

export async function createActivityApi({
  message,
  ttl,
  user_handle,
  authorization,
}) {
  const res = await fetch(`${apiBaseUrl}/api/activities`, {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
      Authorization: authorization,
    },
    body: JSON.stringify({
      user_handle,
      message,
      expires_at: ttlToExpiresAt(ttl),
    }),
  });

  const data = await res.json();

  if (!res.ok) {
    throw new Error(data?.message || "Erro ao criar atividade");
  }

  return data;
}
