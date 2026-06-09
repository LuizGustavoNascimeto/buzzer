import "./HomeFeedPage.css";
import React from "react";
import DesktopNavigation from "../components/DesktopNavigation";
import DesktopSidebar from "../components/DesktopSidebar";
import ActivityFeed from "../components/ActivityFeed";
import ActivityForm from "../components/ActivityForm";
import ReplyForm from "../components/ReplyForm";
import { useAuth } from "../hooks/auth/useAuth";
import { useActivities } from "../hooks/acitivity/useActivities";

export default function HomeFeedPage() {
  const [popped, setPopped] = React.useState(false);
  const [poppedReply, setPoppedReply] = React.useState(false);
  const [replyActivity, setReplyActivity] = React.useState({});

  const { data: activities = [] } = useActivities();
  const { data: user } = useAuth();
  console.log(user)

  return (
    <article>
      <DesktopNavigation user={user} active={"home"} setPopped={setPopped} />
      <div className="content">
        <ActivityForm popped={popped} setPopped={setPopped} />
        <ReplyForm
          activity={replyActivity}
          popped={poppedReply}
          setPopped={setPoppedReply}
          activities={activities}
        />
        <ActivityFeed
          title="Home"
          setReplyActivity={setReplyActivity}
          setPopped={setPoppedReply}
          activities={activities}
        />
      </div>
      <DesktopSidebar handle={user} />
    </article>
  );
}
