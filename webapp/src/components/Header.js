import React from "react";

import { useNavigate } from "react-router-dom";
import Button from "@material-ui/core/Button";
import { logoutUser } from "../services/pollService";
const Header = (props) => {
  const { token, clearToken } = props;
  const navigate = useNavigate();

  const redirectToSignup = () => {
    navigate("/signup");
  };

  const logout = () => {
    logoutUser();
    clearToken();
  };

  const redirectToNew = () => {
    navigate("/new");
  };

  const redirectToLogin = () => {
    navigate("/login");
  };

  return (
    <>
      <div
        style={{
          position: "sticky",
          top: 0,
          padding: "10px 30px",
          height: "40px",
          background: "black",
          color: "white",
        }}
      >
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
          }}
        >
          <>
            <div className="title">polling-app</div>{" "}
            {token ? (
              <Button variant="contained" color="primary" onClick={logout}>
                Logout
              </Button>
            ) : (
              <>
                <Button
                  variant="contained"
                  color="primary"
                  onClick={redirectToLogin}
                >
                  Login
                </Button>
                <Button
                  variant="contained"
                  color="primary"
                  onClick={redirectToSignup}
                >
                  Signup
                </Button>
              </>
            )}
            {token && (
              <button onClick={redirectToNew} style={{ height: "30px" }}>
                Create New poll
              </button>
            )}
          </>
        </div>
      </div>
      <div className="card-header"></div>
      <h1 className="heading">You have {props.numPolls} Polls</h1>
    </>
  );
};
export default Header;
