const apiBaseUrl = process.env.REACT_APP_BACKEND_URL || "http://localhost:8080";
function getAuthHeaders(token) {
  return {
    Accept: "application/json",
    "Content-Type": "application/json",
    Authorization: token ? `Bearer ${token}` : "",
  };
}

export async function fetchMessages(group_id, token) {
  const res = await fetch(`${apiBaseUrl}/api/messages/${group_id}`, {
    method: "GET",
    headers: getAuthHeaders(token),
  });
  const data = await res.json();
  if (!res.ok) {
    throw new Error(data?.message || "Erro ao buscar mensagens");
  }
  return data;
}
export async function createMessageApi(message, token) {
  const res = await fetch(`${apiBaseUrl}/api/messages`, {
    method: "POST",
    headers: getAuthHeaders(token),
    body: JSON.stringify(message),
  });

  const data = await res.json();

  if (!res.ok) {
    throw new Error(data?.message || "Erro ao criar mensagem");
  }

  return data;
}

// type CreateMessageRequest struct {
// 	GroupID        *string `json:"message_group_uuid"`
// 	SenderHandle   string  `json:"sender_handle" binding:"required"`
// 	ReceiverHandle *string `json:"receiver_handle"`
// 	Content        string  `json:"message" binding:"required"`
// }


  //     let json = { 'message': message }
  //     if (params.handle) {
  //       json.handle = params.handle
  //     } else {
  //       json.message_group_uuid = params.message_group_uuid
  //     }

  //     const res = await fetch(backend_url, {
  //       method: "POST",
  //       headers: {
  //         'Authorization': `Bearer ${localStorage.getItem("access_token")}`,
  //         'Accept': 'application/json',
  //         'Content-Type': 'application/json'
  //       },
  //       body: JSON.stringify(json)
  //     });
  //     let data = await res.json();
  //     if (res.status === 200) {
  //       console.log('data:',data)
  //       if (data.message_group_uuid) {
  //         console.log('redirect to message group')
  //         window.location.href = `/messages/${data.message_group_uuid}`
  //       } else {
  //         props.setMessages(current => [...current,data]);
  //       }
  //     } else {
  //       console.log(res)
  //     }
  //   } catch (err) {
  //     console.log(err);
  //   }
  // }