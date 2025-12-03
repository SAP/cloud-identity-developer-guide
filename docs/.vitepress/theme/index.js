import {h, nextTick, watch} from "vue";
import DefaultTheme from "vitepress/theme";
import {useData} from "vitepress";
import {createMermaidRenderer} from "vitepress-mermaid-renderer";
import './custom.css'

export default {
    extends: DefaultTheme,
    Layout: () => {
        const {isDark} = useData();

        const initMermaid = () => {
            const mermaidRenderer = createMermaidRenderer({
                theme: isDark.value ? "dark" : "base"
            });
            mermaidRenderer.setToolbar({
                showLanguageLabel: false,
                desktop: {
                    copyCode: "enabled",
                    toggleFullscreen: "enabled",
                    resetView: "enabled",
                    zoomOut: "enabled",
                    zoomIn: "enabled",
                    zoomLevel: "enabled",
                },
                fullscreen: {
                    copyCode: "disabled",
                    toggleFullscreen: "enabled",
                    resetView: "disabled",
                    zoomLevel: "disabled",
                },
            });
        };

        // initial mermaid setup
        nextTick(() => initMermaid());

        // on theme change, re-render mermaid charts
        watch(
            () => isDark.value,
            () => {
                initMermaid();
            },
        );

        return h(DefaultTheme.Layout);
    },
};