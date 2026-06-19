import React, { useEffect } from "react";
import "./MessageGroupPage.css";

import { useParams } from "react-router-dom";
import DesktopNavigation from "../components/DesktopNavigation";
import MessagesFeed from "../components/MessageFeed";
import MessageGroupFeed from "../components/MessageGroupFeed";
import { useAuth } from "../hooks/auth/useAuth";
import { useMessageGroup } from "../hooks/messageGroups/useMessageGroups";
import { useMessage } from "../hooks/messages/useMessages";

export default function MessageGroupPage() {
  const [popped, setPopped] = React.useState([]);

   const { group_id } = useParams();

  const { data: user, isLoading: userLoading } = useAuth();

  const { data: messageGroups = [] } =
    useMessageGroup(user?.handle, user?.token);

  
  const { data: messages = [] } = useMessage(group_id, user?.token);
  
  useEffect(() => {
    console.log(messages)
  }, [messages])
  
  

  if (userLoading) {
    return (
      <article>
        <div>Carregando usuário...</div>
      </article>
    );
  }

  if (!user) {
    return (
      <article>
        <div>Usuário não autenticado</div>
      </article>
    );
  }
  return (
    <article>
      <DesktopNavigation
        user={user}
        active={"messages"}
        setPopped={setPopped}
      />
      <section className="message_groups">
        <MessageGroupFeed message_groups={messageGroups} />
      </section>
      <div className="content messages">
        <MessagesFeed messages={messages} />
        {/* <MessageForm setMessages={setMessages} /> */}
      </div>
    </article>
  );
}
