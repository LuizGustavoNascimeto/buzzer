import React from "react";
import { useParams } from "react-router-dom";
import "./MessageGroupPage.css";

import DesktopNavigation from "../components/DesktopNavigation";
import MessagesFeed from "../components/MessageFeed";
import MessagesForm from "../components/MessageForm";
import MessageGroupFeed from "../components/MessageGroupFeed";
import { useAuth } from "../hooks/auth/useAuth";
import { useMessageGroup } from "../hooks/messageGroups/useMessageGroups";
import { useCreateMessages } from "../hooks/messages/useCreateMessage";
//import checkAuth from "../lib/CheckAuth";

export default function MessageGroupNewPage() {
  const [otherUser, setOtherUser] = React.useState([]);
  const [message, setMessage] = React.useState([]);
  const [popped, setPopped] = React.useState([]);
  const dataFetchedRef = React.useRef(false);

  const { handle: other_handle } = useParams();

  const { data: user, isLoading: userLoading } = useAuth();
  const { data: messageGroups = [], isLoading: groupsLoading } =
    useMessageGroup(user?.handle, user?.token);
 // const { mutateAsync: createMessage } = useCreateMessages(user?.token);

  // type CreateMessageRequest struct {
  // 	GroupID        *string `json:"message_group_uuid"`
  // 	SenderHandle   string  `json:"sender_handle" binding:"required"`
  // 	ReceiverHandle *string `json:"receiver_handle"`
  // 	Content        string  `json:"message" binding:"required"`
  // }

  const loadUserData = async () => {
    try {
      const backend_url = `${process.env.REACT_APP_BACKEND_URL}/api/users/findByHandle/${other_handle}`;
      const res = await fetch(backend_url, {
        method: "GET",
      });
      let resJson = await res.json();
      if (res.status === 200) {
        setOtherUser(resJson);
      } else {
        console.log(res);
      }
    } catch (err) {
      console.log(err);
    }
  };

  React.useEffect(() => {
    //prevents double call
    if (dataFetchedRef.current) return;
    if (userLoading || !user?.token) return;
    dataFetchedRef.current = true;
    loadUserData();
  }, [userLoading, user]);

  if (userLoading) {
    return <div>Carregando...</div>; // ou um spinner
  }

  return (
    <article>
      <DesktopNavigation active={"messages"} setPopped={setPopped} />
      <section className="message_groups">
        <MessageGroupFeed
          otherUser={otherUser}
          message_groups={messageGroups}
        />
      </section>
      <div className="content messages">
        <MessagesFeed messages={message} />
        <MessagesForm setMessages={setMessage} other_handle={other_handle} />
      </div>
    </article>
  );
}
