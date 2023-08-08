import * as React from "react";

const onSubmit: React.FormEventHandler<HTMLFormElement> = (e) => {
  e.preventDefault();

  const formData = new FormData(e.target as HTMLFormElement);

  fetch("/api/auth", {
    method: "POST",
    body: JSON.stringify(Object.fromEntries(formData.entries())),
    credentials: "include"
  }).then((resp) => {
    console.log(resp);
  });
};

export function LoginBar() {
  return (
    <form onSubmit={onSubmit} method="post" action={"/api/auth"}>
      <input name="username"></input>
      <input name="password"></input>
      <input type="submit"></input>
    </form>
  );
}
