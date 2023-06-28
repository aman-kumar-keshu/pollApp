import axios from "axios";

const ENDPOINT = `http://localhost:8080`;

export const createPoll = async (poll) => {
  console.log("Creating a new poll and saving it to db");

  const res = await axios.post(`${ENDPOINT}/poll`, poll);
  return res.data;
};

export const fetchPolls = async () => {
  console.log("Making call to backend to fetch polls");
  const res = await axios.get(`${ENDPOINT}/polls`);
  return res.data;
};

export const updatePolls = async (id, updatedPoll) => {
  console.log("Updating poll to DB");
  await axios.put(`${ENDPOINT}/poll/${id}`, updatedPoll);
};

export const deletePoll = async (id) => {
  console.log("Deteting the poll");
  await axios.delete(`${ENDPOINT}/poll/${id}`);
};

export const createUser = async (user) => {
  console.log("Creating a user");
  const res = await axios.post(`${ENDPOINT}/users/signup`, user);

  localStorage.setItem("token", res);
  return res;
};

export const loginUser = async (user) => {
  console.log("Login User a user");
  await axios.post(`${ENDPOINT}/users/login`, user);

  // await axios.post(session_url, {}, {
  //   auth: {
  //     username: uname,
  //     password: pass
  //   }
  // });
};
