import "./MessageFeed.css";
import MessageItem from "./MessageItem";
import { useEffect, useRef } from "react";

export default function MessageFeed(props) {
  const collectionRef = useRef(null);

  useEffect(() => {
    if (collectionRef.current) {
      collectionRef.current.scrollTop = collectionRef.current.scrollHeight;
    }
  }, [props.messages]);

  return (
    <div className="message_feed">
      <div className="message_feed_heading">
        <div className="title">Messages</div>
      </div>
      <div className="message_feed_collection" ref={collectionRef}>
        {props.messages.map((message) => {
          return <MessageItem key={message.uuid} message={message} />;
        })}
      </div>
    </div>
  );
}
