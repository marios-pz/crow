import AuthField from "./_components/AuthField/AuthField";

import { useState } from "react";

import { userRegister, userLogin } from "@/utils/auth";

import styles from "./AuthForm.module.scss";

const AuthForm = ({ method }) => {
  let [post, setPost] = useState({
    name: "",
    password: "",
  });

  let [errorMessage, setErrorMessage] = useState("");

  function handleInput(event) {
    setPost({ ...post, [event.target.id]: event.target.value });
  }

  async function handleSubmit(event) {
    event.preventDefault();
    try {
      let response;
      if (method) {
        response = await userLogin(post);
      } else {
        response = await userRegister(post);
      }
      if (response) {
        return (location.href = "/home");
      }
    } catch (error) {
      if (error.response && error.response.data && error.response.data.error) {
        setErrorMessage(`${error.response.data.error}!`);
        return console.error(`Failed to auth user! [${error.message}]`);
      } else {
        return console.error(`Service unavailable! [${error.message}]`);
      }
    }
  }

  return (
    <form className={styles.auth_form} onSubmit={handleSubmit}>
      <AuthField type={1} handler={handleInput} />
      <AuthField type={0} handler={handleInput} />
      <span className={styles.auth_error}>{errorMessage}</span>
      <input
        className={styles.auth_submit}
        type="submit"
        value={method ? "Enter" : "Join"}
      />
    </form>
  );
};

export default AuthForm;
