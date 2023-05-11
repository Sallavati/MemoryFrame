import { createRouter, createWebHistory } from "vue-router";
import Fotoupload from "../views/fotoupload.vue";
import DiashowOnlyPic from "../views/diashow_onlyPic.vue";
import DiashowWithInfo from "../views/diashow_with_Information.vue";
import Settings from "../views/settings.vue";
import Browse from "../views/browse_actions.vue";
import Landing from "../views/landing.vue";
import UploadToExtern from "../views/uploadToExtern.vue";//uploadToExtern

//component: () => import("../views/fetch_test.vue")

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/diashow",
      name: "diashow",
      component: DiashowOnlyPic
    },
    {
      path: "/diashow2",
      name: "diashow2",
      component: DiashowWithInfo
    },
    {
      path: "/settings",
      name: "settings",
      component: Settings
    },
    {
      path: "/",
      name: "landing",
      component: Landing
    },
    {
      path: "/upload",
      name: "fotoupload",
      component: Fotoupload
    },
    {
      path: "/browse",
      name: "browse",
      component: Browse
    },
    {
      path: "/uploadToExtern/:folder",
      name: "Upload To Extern",
      component: UploadToExtern
    }
  ],
});

export default router;
