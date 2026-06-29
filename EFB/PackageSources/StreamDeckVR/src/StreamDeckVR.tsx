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
  private labelsData: { buttons: { row: number; col: number; label: string }[] } = { buttons: [] };
  protected defaultView = "StreamDeckXL";

  protected registerViews(): void {
    this.appViewService.registerPage("StreamDeckXL", () => (
      <StreamDeckXL appViewService={this.appViewService} title="StreamDeck VR"  labelsData={this.labelsData}/>
    ));
    this.defaultView = "StreamDeckXL";
  }

  public render(): VNode {
    return <div class="streamdeckvr-efb">{super.render()}</div>;
  }

  private socket: WebSocket | null = null;

  public async onOpen(): Promise<void> {
    this.socket = new WebSocket("ws://localhost:8080/streamdeckvr");
    console.log("connected")
    this.socket.onmessage = (event) => {
      try {
        const json = JSON.parse(event.data);
        console.log("Received JSON:", json);
        this.labelsData = json;
        this.appViewService.update
      } catch (err) {
        console.error("Failed to parse JSON:", err);
      }
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
