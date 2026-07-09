# HoyoLab Cookie Extraction Guide

[Lifted from here](https://gist.github.com/torikushiii/59eff33fc8ea89dbc0b2e7652db9d3fd) thank you [torikushiii](https://gist.github.com/torikushiii)!

## Overview

This guide provides one method to extract cookies from your HoyoLab profile. These cookies can be used for various purposes, such as automating tasks with scripts. Follow the steps carefully to ensure you obtain the required cookies correctly.

> [WARNING]
> Please login through incognito mode and proceed with the steps below to ensure that you're getting the correct cookies.

## Method 1: Using Network Tab

1. Go to your [HoyoLab profile](https://www.hoyolab.com/accountCenter/postList) and log in with your Hoyoverse account.
2. Open the browser console by pressing `F12`.
3. Navigate to the **Network** tab.

   ![Network Tab](https://gist.github.com/assets/21153445/f3c90ee3-e711-4aea-a1d2-94abd8824c01)
   
4. Search for `getGameRecordCard` in the network requests and click on the result. (If you didn't see anything, refresh the page while keeping the browser console open)

   ![getGameRecordCard Request](https://gist.github.com/assets/21153445/4da91d07-59de-4af0-9471-d1f9e000f61f)
   
5. Go to the **Headers** tab and scroll down to find **Request Headers**. Copy all the cookie values.

   ![Request Headers](https://gist.github.com/assets/21153445/0165a481-682e-411b-ba88-a8af17cd6f71)
6. Copy that entire thing, it should look something like:
```
mi18nLang=xxxx; _MHYUUID=xxxx; HYV_LOGIN_PLATFORM_OPTIONAL_AGREEMENT=xxxx; e_nap_token=xxxx; cookie_token_v2=xxxx; account_mid_v2=xxxx; account_id_v2=xxxx; ltoken_v2=xxxx; ltmid_v2=xxxx; ltuid_v2=xxxx; HYV_LOGIN_PLATFORM_LOAD_TIMEOUT=xxxx; DEVICEFP_SEED_ID=xxxx; DEVICEFP_SEED_TIME=xxxx; DEVICEFP=xxxx; HYV_LOGIN_PLATFORM_TRACKING_MAP=xxxx; HYV_LOGIN_PLATFORM_LIFECYCLE_ID=xxxx
```

Same thing is similar to firefox, look into cookies within headers, look at the `GET` not the `OPTIONS`.


