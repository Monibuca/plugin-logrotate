<template>
    <mu-container>
        <mu-tabs
            :value.sync="active1"
            inverse
            center
        >
            <mu-tab>日志文件</mu-tab>
            <mu-tab>日志跟踪</mu-tab>
            <mu-tab>日志查询</mu-tab>
        </mu-tabs>
        <div class="tabpanel" style="display:flex" v-if="active1 === 0">
            <mu-card v-for="item in logFiles" :key="item.Name">
                <mu-card-title :title="item.Name" :sub-title="unitFormat(item.Size)"></mu-card-title>
                <mu-card-actions>
                    <mu-button flat :href="'/logrotate/open?file='+item.Name" target="_blank">打开</mu-button>
                    <mu-button flat :href="'/logrotate/download?file='+item.Name" target="_blank">下载</mu-button>
                </mu-card-actions>
            </mu-card>
        </div>
        <div class="tabpanel" v-if="active1 === 1">
            <div>
                <mu-switch v-model="autoScroll" label="自动滚动" />
            </div>
            <div ref="logContainer" class="log-container">
                <pre><template v-for="item in logs">{{item+"\n"}}</template></pre>
            </div>
        </div>
        <div class="tabpanel" v-if="active1 === 2">
            <mu-text-field @change="onSearch"></mu-text-field>
            <pre>{{result}}</pre>
        </div>
    </mu-container>
</template>

<script>
let logsES = null;


export default {
    data() {
        return {
            autoScroll: true,
            logs: [],
            logFiles: [],
            result: "",
            active1: 0
        };
    },
    mounted() {
        logsES = new EventSource("/logrotate/tail");
        logsES.onmessage = evt => {
            if (!evt.data) return;
            this.logs.push(evt.data);
        };
        window.ajax.getJSON("/logrotate/list").then(x => (this.logFiles = x));
    },
    deactivated() {
        logsES.close();
    },
    methods: {
        onSearch(value) {
            window.ajax
                .get("/logrotate/find?query=" + value)
                .then(x => (this.result = x));
        }
    },
    updated() {
        if (this.autoScroll) {
            this.$refs.logContainer.scrollTop = this.$refs.logContainer.offsetHeight;
        }
    }
};
</script>

<style>
.log-container {
    overflow-y: auto;
    max-height: 500px;
}
.tabpanel {
    padding: 0 20px;
}
</style>