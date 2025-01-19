import axios from "../services/axios";

import { authStore } from "../stores/authStore";

export const uploadImage = async (data) => {
  try {
    const formData = new FormData();
    formData.append("image", data);
    // Upload the image
    const resUpload = await axios.post("/upload", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
        Authorization: `Bearer ${authStore.user.data.userID}`,
      },
    });

    return resUpload.data.image;
  } catch (error) {
    throw new Error(error);
  }
};
