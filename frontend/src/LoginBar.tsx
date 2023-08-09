import * as React from "react";
import Cookies from "js-cookie";

interface LoginBarProps {
  username: string | null
}

export function LoginBar(props: LoginBarProps) {
  const [errors, setErrors] = React.useState<string | null>(null);

  const onSubmit: React.FormEventHandler<HTMLFormElement> = (e) => {
    e.preventDefault();

    const form = e.target as HTMLFormElement;
    const formData = new FormData(form);

    fetch(form.action, {
      method: form.method,
      body: JSON.stringify(Object.fromEntries(formData.entries())),
      credentials: "include",
      headers: {},
    }).then(async (resp) => {
      if (resp.ok) {
        location.reload();
        return;
      }
      const x = await resp.text();
      setErrors(x);
    });
  };


  if (props.username !== null) {
    return (
      <>
        <p>The currently logged in user is {props.username}</p>
        <button
          onClick={() => {
            Cookies.remove("user_session", {path: "/api"});
            location.reload();
          }}
        >
          Loggout
        </button>
      </>
    );
  }

  return (
    <>
      <form onSubmit={onSubmit} method="post" action={"/api/auth"}>
        <input name="username"></input>
        <input name="password"></input>
        <input type="submit"></input>
      </form>
      {errors === null ? <></> : <p>{errors}</p>}
    </>
  );
}
