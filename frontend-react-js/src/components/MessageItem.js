import "./MessageItem.css";
import { Link } from "react-router-dom";
import { DateTime } from "luxon";
import { useAuth } from "../hooks/auth/useAuth";

export default function MessageItem(props) {
  const format_time_created_at = (value) => {
    // format: 2050-11-20 18:32:47 +0000
    const created = DateTime.fromISO(value);
    const now = DateTime.now();
    const diff_mins = now.diff(created, "minutes").toObject().minutes;
    const diff_hours = now.diff(created, "hours").toObject().hours;
    if (diff_hours > 24.0) {
      return created.toFormat("LLL L");
    } else if (diff_hours < 24.0 && diff_hours > 1.0) {
      return `${Math.floor(diff_hours)}h`;
    } else if (diff_hours < 1.0) {
      return `${Math.round(diff_mins)}m`;
    }
  };

  const { data: user, isLoading: userLoading } = useAuth();

  // enquanto não sabemos quem está logado, tratamos como "não é minha mensagem"
  // pra não piscar o card pro lado errado quando o auth carregar
  const is_own = !userLoading && user?.handle === props.message.user_handle;

  return (
    <Link
      className={`message_item${is_own ? " is_own" : ""}`}
      to={`/messages/@` + props.message.user_handle}
    >
      <div className="message_header">
        <div className="message_avatar"></div>
        <div className="message_identity">
          <span className="display_name">
            {props.message.user_display_name}
          </span>
          <span className="handle">@{props.message.user_handle}</span>
        </div>
        {/* message_identity */}
      </div>
      {/* message_header */}
      <div className="message_content">
        <div className="message">{props.message.message}</div>
        <div className="created_at" title={props.message.created_at}>
          <span className="ago">
            {format_time_created_at(props.message.created_at)}
          </span>
        </div>
        {/* created_at */}
      </div>
      {/* message_content */}
    </Link>
  );
}