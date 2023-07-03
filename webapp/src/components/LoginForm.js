import React, { useState } from "react";
import { makeStyles } from "@material-ui/core";
import TextField from "@material-ui/core/TextField";
import Button from "@material-ui/core/Button";
import { loginUser } from "../services/pollService";
import { useNavigate } from "react-router-dom";

const useStyles = makeStyles((theme) => ({
  root: {
    display: "flex",
    flexDirection: "column",
    justifyContent: "center",
    alignItems: "center",
    padding: theme.spacing(2),

    "& .MuiTextField-root": {
      margin: theme.spacing(1),
      width: "300px",
    },
    "& .MuiButtonBase-root": {
      margin: theme.spacing(2),
    },
  },
}));

const LoginForm = () => {
  const navigator = useNavigate();
  const classes = useStyles();
  // create state variables for each input
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState(null);

  const goBack = () => {
    navigator("/");
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const res = await loginUser({ email, password });

      console.log("Logged in successfully", res);
      alert("Logged in Successfully");

      navigator("/");

      navigator("/");
    } catch (error) {
      alert(error.response.data.message);
      navigator("/login");
    }
  };
  return (
    <form className={classes.root} onSubmit={handleSubmit}>
      <TextField
        label="Email"
        variant="filled"
        type="email"
        required
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />
      <TextField
        label="Password"
        variant="filled"
        type="password"
        required
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      <div>
        <Button variant="contained" onClick={goBack}>
          Go Back
        </Button>
        <Button
          type="submit"
          onClick={handleSubmit}
          variant="contained"
          color="primary"
        >
          LOGIN
        </Button>
      </div>
    </form>
  );
};

export default LoginForm;
