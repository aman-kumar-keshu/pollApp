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
  console.log(res.request.status);
  console.log(res.data.token);
  if (res.request.status === 200) {
    localStorage.setItem("token", res.data.token);
  }
};

export const loginUser = async (user) => {
  console.log("Login User a user");
  try {
    const res = await axios.post(`${ENDPOINT}/users/login`, user);
    console.log(res);
    if (res.request.status === 202) {
      console.log("Logged in successfully", res.data.token);
      localStorage.setItem("token", res.data.token);
    }
  } catch (error) {
    
  }
};

export const logoutUser = async () => {
  console.log("logout user");
  localStorage.removeItem("token");
};
