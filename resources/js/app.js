require("./bootstrap");

// Import modules...
import Vue from 'vue';
import { createInertiaApp } from '@inertiajs/inertia-vue';

import { InertiaProgress } from "@inertiajs/progress";
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
dayjs.extend(relativeTime);

Vue.use(InertiaProgress);

createInertiaApp({
    resolve: name => require(`./Pages/${name}`),
    setup({ el, App, props }) {
        new Vue({
            render: h => h(App, props),
        }).$mount(el)
    },
});

InertiaProgress.init({ color: "#4B5563" });
