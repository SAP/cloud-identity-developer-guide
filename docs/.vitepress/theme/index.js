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
                theme: isDark.value ? "dark" : "base",
                // Make every diagram type honor the parent container width by default
                // so the diagram is fully visible at scale=1 and the user can zoom in
                // via the toolbar if needed.
                flowchart: { useMaxWidth: true },
                sequence: { useMaxWidth: true },
                gantt: { useMaxWidth: true },
                journey: { useMaxWidth: true },
                class: { useMaxWidth: true },
                state: { useMaxWidth: true },
                er: { useMaxWidth: true },
                pie: { useMaxWidth: true },
                requirement: { useMaxWidth: true },
                mindmap: { useMaxWidth: true },
                timeline: { useMaxWidth: true },
                gitGraph: { useMaxWidth: true },
                c4: { useMaxWidth: true },
                quadrantChart: { useMaxWidth: true },
                xyChart: { useMaxWidth: true },
                sankey: { useMaxWidth: true },
                block: { useMaxWidth: true },
                packet: { useMaxWidth: true },
                architecture: { useMaxWidth: true }
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