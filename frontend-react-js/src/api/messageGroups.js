const apiBaseUrl =
  process.env.REACT_APP_BACKEND_URL || "http://localhost:8080";

function getAuthHeaders(token) {
  return {
    Accept: "application/json",
    "Content-Type": "application/json",
    Authorization: token ? `Bearer ${token}` : "",
  };
}

export async function fetchMessageGroup(handle, token) {
  const res = await fetch(
    `${apiBaseUrl}/api/message_groups/${handle}`,
    {
      method: "GET",
      headers: getAuthHeaders(token),
    }
  );

  const data = await res.json();

  if (!res.ok) {
    throw new Error(data?.message || "Erro ao buscar grupo de mensagens");
  }

  return data;
}