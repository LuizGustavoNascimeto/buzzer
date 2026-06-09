import "./ActivityForm.css";
import React from "react";
import process from "process";
import { ReactComponent as BombIcon } from "./svg/bomb.svg";
import { useCreateActivity } from "../hooks/useCreateActivity";
import { useAuth } from "../hooks/useAuth";

export default function ActivityForm(props) {
  const [count, setCount] = React.useState(0);
  const [message, setMessage] = React.useState("");
  const [ttl, setTtl] = React.useState("7-days");

  const classes = [];
  classes.push("count");
  if (240 - count < 0) {
    classes.push("err");
  }

  const { mutate: createActivity } = useCreateActivity();
  const { user, token } = useAuth();

  const onSubmit = async (event) => {
    event.preventDefault();
    try {
      await createActivity({
        message,
        ttl,
        user_handle: user.handle,
        authorization: token,
      });

      setCount(0);
      setMessage("");
      setTtl("7-days");
      props.setPopped(false);
    } catch (err) {
      console.error(err);
    }
  };

  const textarea_onchange = (event) => {
    setCount(event.target.value.length);
    setMessage(event.target.value);
  };

  const ttl_onchange = (event) => {
    setTtl(event.target.value);
  };

  if (props.popped === true) {
    return (
      <form className="activity_form" onSubmit={onSubmit}>
        <textarea
          type="text"
          placeholder="what would you like to say?"
          value={message}
          onChange={textarea_onchange}
        />
        <div className="submit">
          <div className={classes.join(" ")}>{240 - count}</div>
          <button type="submit">Buzz</button>
          <div className="expires_at_field">
            <BombIcon className="icon" />
            <select value={ttl} onChange={ttl_onchange}>
              <option value="30-days">30 days</option>
              <option value="7-days">7 days</option>
              <option value="3-days">3 days</option>
              <option value="1-day">1 day</option>
              <option value="12-hours">12 hours</option>
              <option value="3-hours">3 hours</option>
              <option value="1-hour">1 hour </option>
            </select>
          </div>
        </div>
      </form>
    );
  }
}
