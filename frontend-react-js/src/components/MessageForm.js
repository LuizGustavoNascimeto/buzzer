import React from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useAuth } from "../hooks/auth/useAuth";
import { useCreateMessages } from "../hooks/messages/useCreateMessage";
import "./MessageForm.css";

export default function MessageForm(props) {
  const [count, setCount] = React.useState(0);
  const [message, setMessage] = React.useState("");

  const navigate = useNavigate();

  const classes = [];
  classes.push("count");
  if (1024 - count < 0) {
    classes.push("err");
  }

  const { data: user } = useAuth();
  const { group_id } = useParams();
  const { mutateAsync: createMessage } = useCreateMessages(
    group_id,
    user.token,
  );
  let redirect_id;
  // o group_id só existe quando a conversa ja existe, enquanto o props.other_handle só existe quando a conversa é nova
  // a dois fluxos aqui, escondidos nessa ideia

  const onsubmit = async (event) => {
    event.preventDefault();
    try {
      const payload = {
        message: message,
        message_group_uuid: group_id,
        sender_handle: user.handle,
        receiver_handle: props.other_handle,
      };

      const data = await createMessage(payload, user.token);
      redirect_id = data.message_group_uuid;

      props.setMessages((current) => [...current, data]);

      setMessage("");
      setTimeout(() => {
        console.log("message state:", message);
      }, 100);
      setCount(0);
    } catch (err) {
      console.log(err);
    } finally {
      if (!group_id) {
        navigate(`/messages/${redirect_id}`);
      }
    }
  };

  const textarea_onchange = (event) => {
    setCount(event.target.value.length);
    setMessage(event.target.value);
  };

  return (
    <form className="message_form" onSubmit={onsubmit}>
      <textarea
        type="text"
        placeholder="send a direct message..."
        value={message}
        onChange={textarea_onchange}
      />
      <div className="submit">
        <div className={classes.join(" ")}>{1024 - count}</div>
        <button type="submit">Message</button>
      </div>
    </form>
  );
}
