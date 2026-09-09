import {
  App,
  AppBootMode,
  AppInstallProps,
  AppSuspendMode,
  AppView,
  AppViewProps,
  Efb,
  RequiredProps,
  TVNode,
} from "@efb/efb-api";
import { FSComponent, VNode } from "@microsoft/msfs-sdk";
import { StreamDeckXL } from "./Components/StreamDeckXL";

import "./StreamDeckVR.scss";

declare const BASE_URL: string;

class StreamDeckVRAppView extends AppView<RequiredProps<AppViewProps, "bus">> {
  protected defaultView = "StreamDeckXL";

  protected registerViews(): void {
    this.appViewService.registerPage("StreamDeckXL", () => (
      <StreamDeckXL appViewService={this.appViewService} bus={this.bus} title="StreamDeck VR"/>
    ));
    this.defaultView = "StreamDeckXL";
  }

  public render(): VNode {
    return <div class="streamdeckvr-efb">{super.render()}</div>;
  }

  private socket: WebSocket | null = null;

  public async onOpen(): Promise<void> {
    console.log("[StreamDeckVR] Attempting to connect to WebSocket at ws://localhost:8081/streamdeckvr...");
    try {
        this.socket = new WebSocket("ws://localhost:8081/streamdeckvr");
        this.socket.onopen = () => {
            console.log("[StreamDeckVR] WebSocket connection established!");
        };
        this.socket.onerror = (error) => {
            console.error("[StreamDeckVR] WebSocket error:", error);
        };
        this.socket.onclose = (event) => {
            console.log(
                `[StreamDeckVR] WebSocket connection closed. Code: ${event.code}, Reason: ${event.reason || "No reason provided"}`
            );
        };
        this.socket.onmessage = (event) => {
            console.log("[StreamDeckVR] Raw message received:", event.data);

            try {
                const json = JSON.parse(event.data);
                console.log("[StreamDeckVR] Parsed JSON:", json);
                this.bus.pub("streamdeck-labels-updated", json);
                this.appViewService.update;
            } catch (err) {
                console.error("[StreamDeckVR] Failed to parse JSON:", err);
                console.error("[StreamDeckVR] Raw data that failed to parse:", event.data);
            }
        };
    } catch (err) {
        console.error("[StreamDeckVR] Failed to initialize WebSocket:", err);
    }
  };


  public async onClose(): Promise<void> {
    if (this.socket) {
      this.socket.close();
      this.socket = null;
      console.log("disconnected")
    }
  }

  //   public async onPause(): Promise<void> {
  //   if (this.socket) {
  //     this.socket.close();
  //     this.socket = null;
  //     console.log("disconnected")
  //   }
  // }
}

class StreamDeckVR extends App {
 
  public get name(): string {
    return StreamDeckVR.name;
  }
 
  public get icon(): string {
    return `${BASE_URL}/Assets/app-icon.svg`;
  }

  /**
   * Optional attribute
   * Allow to choose BootMode between COLD / WARM / HOT
   * Default behavior : AppBootMode.COLD
   *
   * COLD : No dom preloaded in memory
   * WARM : App -> AppView are loaded but not rendered into DOM
   * HOT : App -> AppView -> Pages are rendered and injected into DOM
   */
  public BootMode = AppBootMode.COLD;

  /**
   * Optional attribute
   * Allow to choose SuspendMode between SLEEP / TERMINATE
   * Default behavior : AppSuspendMode.SLEEP
   *
   * SLEEP : Default behavior, does nothing, only hiding and sleeping the app if switching to another one
   * TERMINATE : Hiding the app, then killing it by removing it from DOM (BootMode is checked on next frame to reload it and/or to inject it, see BootMode)
   */
  public SuspendMode = AppSuspendMode.SLEEP;

  /**
   * Optional method
   * Allow to resolve some dependencies, install external data, check an api key, etc...
   * @param _props props used when app has been setted up.
   * @returns Promise<void>
   */
  public async install(_props: AppInstallProps): Promise<void> { 
    Efb.loadCss(`${BASE_URL}/StreamDeckVR.css`);
    return Promise.resolve();
  }

  public render(): TVNode<StreamDeckVRAppView> {
    return <StreamDeckVRAppView bus={this.bus} />;
  }
}

Efb.use(StreamDeckVR);
