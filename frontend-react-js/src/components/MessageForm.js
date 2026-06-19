import "./MessageForm.css";
import React from "react";
import process from "process";
import { useParams } from "react-router-dom";
import { useCreateMessages } from "../hooks/messages/useCreateMessage";
import { useAuth } from "../hooks/auth/useAuth";

export default function MessageForm(props) {
  const [count, setCount] = React.useState(0);
  const [message, setMessage] = React.useState("");
  const params = useParams();

  const classes = [];
  classes.push("count");
  if (1024 - count < 0) {
    classes.push("err");
  }

  const { mutate: createMessage } = useCreateMessages();

  const { data: user } = useAuth();
  const onsubmit = async (event) => {
    event.preventDefault();

    try {
      const payload = {
        message: message,
        user_receiver_handle: params.handle,
      };

      console.log("onsubmit payload", payload);

      const data = await createMessage(payload);

      props.setMessages((current) => [...current, data]);

      setMessage("");
      setCount(0);
    } catch (err) {
      console.log(err);
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
