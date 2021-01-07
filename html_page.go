package golive

var BasePageString = `<!DOCTYPE html>
<html lang="{{ .Lang }}">
  <head>
    <meta charset="UTF-8" />
    <title>{{ .Title }}</title>
    {{ .Head }}
  </head>
  <script type="application/javascript">
    const GO_LIVE_CONNECTED = "go-live-connected";
    const GO_LIVE_COMPONENT_ID = "go-live-component-id";
    const EVENT_LIVE_DOM_COMPONENT_ID_KEY = 'cid'
    const EVENT_LIVE_DOM_INSTRUCTIONS_KEY = 'i'
    const EVENT_LIVE_DOM_TYPE_KEY = 't'
    const EVENT_LIVE_DOM_CONTENT_KEY = 'c'
    const EVENT_LIVE_DOM_ATTR_KEY = 'a'
    const EVENT_LIVE_DOM_SELECTOR_KEY = 's'

    const findLiveInputsFromElement = (el) => {
      return el.querySelectorAll(
        ['*[go-live-input]:not([', GO_LIVE_CONNECTED, '])'].join('')
      );
    };

    const findLiveClicksFromElement = (el) => {
      return el.querySelectorAll(
        ['*[go-live-click]:not([', GO_LIVE_CONNECTED, '])'].join('')
      );
    };

    function getElementChild(element, index) {
      let el = element.firstChild;

      if (el === Node.TEXT_NODE) {
        throw new Error("Element is a text node, without children");
      }

      while (index > 0) {
        if (!el) {
          console.log("Element not found in path", element);
          return;
        }

        el = el.nextSibling;
        index--;
      }

      return el;
    }

    function isElement(o) {
      return typeof HTMLElement === "object"
        ? o instanceof HTMLElement //DOM2
        : o &&
            typeof o === "object" &&
            o.nodeType === 1 &&
            typeof o.nodeName === "string";
    }

    function handleDiffSetAttr(message, el) {
      const { attr } = message;

      if (attr.Name === "value" && el.value) {
        el.value = attr.Value;
      } else {
        el.setAttribute(attr.Name, attr.Value);
      }
    }

    function handleDiffRemoveAttr(message, el) {
      const { attr } = message;

      el.removeAttribute(attr.Name);
    }

    function handleDiffReplace(message, el) {
      const { content } = message;

      const wrapper = document.createElement("div");
      wrapper.innerHTML = content;

      el.parentElement.replaceChild(wrapper.firstChild, el);
    }

    function handleDiffRemove(message, el) {
      el.parentElement.removeChild(el);
    }

    function handleDiffSetInnerHTML(message, el, componentId) {
      const { content } = message;

      if (el.nodeType === Node.TEXT_NODE) {
        el.textContent = content;
        return;
      }

      el.innerHTML = content;

      goLive.connectElement( el);
    }

    function handleDiffAppend(message, el, componentId) {
      const { content } = message;

      const wrapper = document.createElement("div");
      wrapper.innerHTML = content;
      const child = wrapper.firstChild;
      el.appendChild(child);

      goLive.connectElement( el);
    }

    const handleChange = {
      "{{ .Enum.DiffSetAttr }}": handleDiffSetAttr,
      "{{ .Enum.DiffRemoveAttr }}": handleDiffRemoveAttr,
      "{{ .Enum.DiffReplace }}": handleDiffReplace,
      "{{ .Enum.DiffRemove }}": handleDiffRemove,
      "{{ .Enum.DiffSetInnerHTML }}": handleDiffSetInnerHTML,
      "{{ .Enum.DiffAppend }}": handleDiffAppend,
    };

    const createOnceEmitter = () => {
      const handlers = {
      }
      const createHandler = (name, called) => {
        handlers[name] = {
          called,
          cbs: [],
        }

        return handlers[name]
      }

      return {
        on(name, cb) {
          let handler = handlers[name];

          if (!handler) {
            handler = createHandler(name, false);
          }

          handler.cbs.push(cb);
        },
        emit(name, ...attrs) {
          const handler = handlers[name];

          if (!handler) {
            createHandler(name, true);
            return;
          }

          for (const cb of handler.cbs) {
            cb();
          }
        }
      }
    }

    const getComponentIdFromElement = (element) => {
      const attr = element.getAttribute("go-live-component-id")
      if ( attr ) {
        return attr
      }

      if(element.parentElement) {
        return getComponentIdFromElement(element.parentElement)
      }

      return undefined
    }

    const goLive = {

      server: new WebSocket(["ws://", window.location.host, "/ws"].join("")),

      handlers: [],
      onceHandlers: {},

      once: createOnceEmitter(),

      getLiveComponent(id) {
        return document.querySelector(
          ["*[",GO_LIVE_COMPONENT_ID, "=", id, "]"].join("")
        );
      },

      on(name, handler) {
        const newSize = this.handlers.push({
          name,
          handler,
        });
        return newSize - 1;
      },

      findHandler(name) {
        return this.handlers.filter((i) => i.name === name);
      },

      emit(name, message) {
        for (const handler of this.findHandler(name)) {
          handler.handler(message);
        }
      },

      off(index) {
        this.handlers.splice(index, 1);
      },

      send(message) {
        goLive.server.send(
            JSON.stringify(message))
      },

      connectChildren(viewElement) {
        const liveChildren = viewElement.querySelectorAll(
          "*[" + GO_LIVE_COMPONENT_ID + "]"
        );

        liveChildren.forEach((child) => {
          const componentId = child.getAttribute(GO_LIVE_COMPONENT_ID);
          this.connectElement( child);
        });
      },

      connectElement(viewElement) {
        if (typeof viewElement === "string") {
          return;
        }

        if (!isElement(viewElement)) {
          return;
        }

        const liveInputs = findLiveInputsFromElement(viewElement);
        const clickElements = findLiveClicksFromElement(viewElement);


        clickElements.forEach(function (element) {
          const componentId = getComponentIdFromElement(element)
          element.addEventListener("click", function (_) {
            goLive.send({
                name: "{{ .Enum.EventLiveMethod }}",
                component_id: componentId,
                method_name: element.getAttribute("go-live-click"),
                value: String(element.value),
            })
          });
          element.setAttribute(GO_LIVE_CONNECTED, true);
        });

        liveInputs.forEach(function (element) {
          const type = element.getAttribute("type");
          const componentId = getComponentIdFromElement(element)

          element.addEventListener("input", function (_) {
            let value = element.value;

            if (type === "checkbox") {
              value = element.checked;
            }

            goLive.send({
                name: "{{ .Enum.EventLiveInput }}",
                component_id: componentId,
                key: element.getAttribute("go-live-input"),
                value: String(value),
            })

          });
          element.setAttribute(GO_LIVE_CONNECTED, true);
        });
      },

                        if (!el) {
                            console.warn("Path not found with selector", path)
                            return
                        }

                        if (el.nodeType === Node.TEXT_NODE) {
                            el.textContent = content
                            return
                        }

                        el.innerHTML = content

                        // Add listeners to new elements
                        goLive.connectElement(scopeId, el)
                    },
                    '{{ .Enum.DiffAppend }}': (message) => {
                        const {
                            content,
                            path
                        } = message

                        const el = getElementByIndexPath(path, viewElement)

                        if (!el) {
                            console.warn("Path not found with selector", path)
                            return
                        }

                        const wrapper = document.createElement('div');
                        wrapper.innerHTML = content;

                        if (content.trim().length > 0) {
                            el.appendChild(wrapper.firstChild)
                        }

                        goLive.connectElement(scopeId, el.lastChild)
                    }
                }

                if (viewElement.getAttribute("go-live-component-id") === message.component_id) {
                    handleChange[message.type.toLowerCase()](message)
                }
            })
        },

        connect(id) {
            goLive.connectElement(id, goLive.getLiveComponent(id))
        },
    }

    goLive.server.onmessage = (rawMessage) => {
        const message = JSON.parse(rawMessage.data)
        goLive.emit(message.name, message)
    }

    goLive.server.onopen = () => {
        goLive.emitOnce('WS_CONNECTION_OPEN')
    }

</script>
<body>
{{ .Body }}
</body>
</html>
`
