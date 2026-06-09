import "./HomeFeedPage.css";
import React, { useEffect } from "react";

import { getCurrentUser, fetchUserAttributes } from "aws-amplify/auth";
import DesktopNavigation from "../components/DesktopNavigation";
import DesktopSidebar from "../components/DesktopSidebar";
import ActivityFeed from "../components/ActivityFeed";
import ActivityForm from "../components/ActivityForm";
import ReplyForm from "../components/ReplyForm";
import { useAuth } from "../hooks/useAuth";
import { useActivities } from "../hooks/useActivities";

export default function HomeFeedPage() {
  const [popped, setPopped] = React.useState(false);
  const [poppedReply, setPoppedReply] = React.useState(false);
  const [replyActivity, setReplyActivity] = React.useState({});
  const dataFetchedRef = React.useRef(false);
  const { user, error } = useAuth();

  const { data: activities = [], isLoading } = useActivities();

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
      <DesktopSidebar user={user} />
    </article>
  );
}
