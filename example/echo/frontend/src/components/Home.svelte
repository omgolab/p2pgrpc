<script lang="ts">
  import {createPromiseClient} from '@connectrpc/connect';
  import {createConnectTransport} from '@connectrpc/connect-web';
  import {DemoService} from 'demo_services/v1/demo_services_connect';
  import type {ListPromptsResponse,StreamMovieNamesResponse} from 'demo_services/v1/demo_services_pb';

  let prompts: ListPromptsResponse | null = null;
  let movieNames: StreamMovieNamesResponse[] = [];

  const client = createPromiseClient(
    DemoService,
    createConnectTransport({
      baseUrl: 'http://localhost:8015',
      // baseUrl: 'https://gw.stage.astronlab.com',
      useBinaryFormat: false,
    })
  );

  const handleButtonClick = async () => {
    const response = await client.listPrompts({});
    prompts = response;
    console.log(response);
  };

  const handlePostClick = async () => {
    console.log('saving');
    const response = await client.savePrompt({
      text: 'Hello from Svelte',
    });
    console.log(response);
  };

  const handleStreamClick = async () => {
    const stream = client.streamMovieNames({});
    for await (const response of stream) {
      console.log(response);
      movieNames=[...movieNames,response];
    }
  };

  const handleClientStreamClick = async () => {
  //  now stream to server from client
    try {
      for (let i = 0; i < 10; i++) {
        const response = await client.streamQuotes({numOfQuotes: 5});
        console.log(response);
      }
    } catch (e) {
      console.log(e);
    }
  };
</script>

<div class="hero min-h-screen bg-base-200">
  <div class="hero-content text-center">
    <div class="max-w-md">
      <h1 class="text-5xl font-bold">Welcome!</h1>
      <p class="py-6">
        Test the connection to the server by fetching the list of prompts.
      </p>
      <button on:click={handleButtonClick} class="btn btn-primary">Fetch</button
      >
      <button on:click={handleStreamClick} class="btn btn-primary"
        >Server Stream</button
      >

      <button on:click={handleClientStreamClick} class="btn btn-primary"
        >Client Stream</button
      >
      <button on:click={handlePostClick} class="btn btn-primary">Save</button>

      <div class="overflow-x-auto mt-10 card bg-base-100 shadow-xl text-center max-h-96 overflow-y-scroll">
        {#if prompts}
          <table class="table">
            <thead>
              <tr class="text-center">
                <th>SL</th>
                <th>ID</th>
                <th>Text</th>
              </tr>
            </thead>
            <tbody class="text-center">
              {#each prompts.prompts as prompt, i}
                <tr>
                  <td>{i + 1}</td>
                  <td>{prompt.id}</td>
                  <td>{prompt.text}</td>
                </tr>
              {/each}
            </tbody>
          </table>
          {/if}
        {#if movieNames.length > 0}
          <table class="table">
            <thead>
              <tr class="text-center">
                <th>SL</th>
                <th>Movie Name</th>
              </tr>
            </thead>
            <tbody class="text-center ">
              {#each movieNames as movieName, i}
                <tr>
                  <td>{i + 1}</td>
                  <td>{movieName.movieName}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        {/if}

      </div>
    </div>
  </div>
</div>
