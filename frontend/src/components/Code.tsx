import { useEffect, useState } from "react";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { vscDarkPlus } from "react-syntax-highlighter/dist/esm/styles/prism";

export function Code({ functionName }: { functionName: string }) {
  const [error, setError] = useState<string>("");
  const [codeText, setCodeText] = useState<string>("");
  useEffect(() => {
    async function loadData() {
      const resp = await fetch(
        `http://localhost:8090/func?functionName=${functionName}`
      );
      if (!resp.ok) {
        setError("Could not get function code");
        return;
      }

      const data = await resp.text();
      setCodeText(data);
    }
    loadData();
  }, [functionName]);
  if (error != "") {
    return <p>{error}</p>;
  }
  return (
    <SyntaxHighlighter language="go" style={vscDarkPlus}>
      {codeText}
    </SyntaxHighlighter>
  );
}
