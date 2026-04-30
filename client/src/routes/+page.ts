import { browser } from "$app/environment";
import { redirect } from "@sveltejs/kit";

export const ssr = false;

export const load = () => {
  if (browser) {
    redirect(302, localStorage.getItem("token") ? "/home" : "/auth");
  }
};
