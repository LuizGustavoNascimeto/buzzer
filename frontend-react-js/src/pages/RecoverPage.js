import "./RecoverPage.css";
import React from "react";
import { ReactComponent as Logo } from "../components/svg/logo.svg";
import { Link } from "react-router-dom";
import { resetPassword, confirmResetPassword } from "aws-amplify/auth";

export default function RecoverPage() {
  // Username is Eamil
  const [username, setUsername] = React.useState("");
  const [password, setPassword] = React.useState("");
  const [passwordAgain, setPasswordAgain] = React.useState("");
  const [code, setCode] = React.useState("");
  const [errors, setErrors] = React.useState("");
  const [formState, setFormState] = React.useState("send_code");

  const onsubmit_send_code = async (event) => {
    event.preventDefault();
    setErrors("");
    try {
      const output = await resetPassword({ username });
      const { nextStep } = output;
      switch (nextStep.resetPasswordStep) {
        case "CONFIRM_RESET_PASSWORD_WITH_CODE":
          const codeDeliveryDetails = nextStep.codeDeliveryDetails;
          console.log(
            `Confirmation code was sent to ${codeDeliveryDetails.deliveryMedium}`,
          );
          // Collect the confirmation code from the user and pass to confirmResetPassword.
          break;
        case "DONE":
          console.log("Successfully reset password.");
          break;
      }
    } catch (error) {
      setErrors(error.message);
    }
    finally {
      setFormState('confirm_code')
    }
  };

  const onsubmit_confirm_code = async (event) => {
    event.preventDefault();
    setErrors("");
    if (password == passwordAgain) {
      try {
        await confirmResetPassword({ username, code, password });
      } catch (error) {
        setErrors(error.message);
      } finally {
        setFormState("success");
      }
    } else {
      setErrors("Passwords do not match");
    }
  };

  const username_onchange = (event) => {
    setUsername(event.target.value);
  };
  const password_onchange = (event) => {
    setPassword(event.target.value);
  };
  const password_again_onchange = (event) => {
    setPasswordAgain(event.target.value);
  };
  const code_onchange = (event) => {
    setCode(event.target.value);
  };
  const send_code = () => {
    return (
      <form className="recover_form" onSubmit={onsubmit_send_code}>
        <h2>Recover your Password</h2>
        <div className="fields">
          <div className="field text_field username">
            <label>Email</label>
            <input type="text" value={username} onChange={username_onchange} />
          </div>
        </div>
        {errors && <div className="errors">{errors}</div>}
        <div className="submit">
          <button type="submit">Send Recovery Code</button>
        </div>
      </form>
    );
  };

  const confirm_code = () => {
    return (
      <form className="recover_form" onSubmit={onsubmit_confirm_code}>
        <h2>Recover your Password</h2>
        <div className="fields">
          <div className="field text_field code">
            <label>Reset Password Code</label>
            <input type="text" value={code} onChange={code_onchange} />
          </div>
          <div className="field text_field password">
            <label>New Password</label>
            <input
              type="password"
              value={password}
              onChange={password_onchange}
            />
          </div>
          <div className="field text_field password_again">
            <label>New Password Again</label>
            <input
              type="password"
              value={passwordAgain}
              onChange={password_again_onchange}
            />
          </div>
        </div>
        {errors}
        <div className="submit">
          <button type="submit">Reset Password</button>
        </div>
      </form>
    );
  };

  const success = () => {
    return (
      <form>
        <p>Your password has been successfully reset!</p>
        <Link to="/signin" className="proceed">
          Proceed to Signin
        </Link>
      </form>
    );
  };

  let form;
  if (formState == "send_code") {
    form = send_code();
  } else if (formState == "confirm_code") {
    form = confirm_code();
  } else if (formState == "success") {
    form = success();
  }

  return (
    <article className="recover-article">
      <div className="recover-info">
        <Logo className="logo" />
      </div>
      <div className="recover-wrapper">{form}</div>
    </article>
  );
}
