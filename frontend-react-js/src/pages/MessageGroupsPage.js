import "./MessageGroupsPage.css";
import React from "react";

import DesktopNavigation from "../components/DesktopNavigation";
import MessageGroupFeed from "../components/MessageGroupFeed";
import { useAuth } from "../hooks/auth/useAuth";
import { useMessageGroup } from "../hooks/messageGroups/useMessageGroups";

export default function MessageGroupsPage() {
  const [popped, setPopped] = React.useState([]);

  const { data: user, isLoading: userLoading } = useAuth();

const { data: messageGroups = [], isLoading: groupsLoading } =
  useMessageGroup(user?.handle, user?.token);

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
        {groupsLoading ? (
          <div>Carregando mensagens...</div>
        ) : (
          <MessageGroupFeed message_groups={messageGroups} />
        )}
      </section>

      <div className="content" />
    </article>
  );
}
