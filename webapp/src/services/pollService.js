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
