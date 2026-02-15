package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chromedp/chromedp"
)

// clickByNodeID resolves a backend DOM node to a JS object and calls .click().
// Uses DOM.resolveNode + Runtime.callFunctionOn which works on React/shadow DOM
// where CSS selectors fail.
func clickByNodeID(ctx context.Context, backendNodeID int64) error {
	return chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			objectID, err := resolveNodeToObject(ctx, backendNodeID)
			if err != nil {
				return err
			}
			callP := map[string]any{
				"objectId":            objectID,
				"functionDeclaration": "function() { this.scrollIntoViewIfNeeded(); this.click(); }",
				"arguments":           []any{},
			}
			if err := chromedp.FromContext(ctx).Target.Execute(ctx, "Runtime.callFunctionOn", callP, nil); err != nil {
				return fmt.Errorf("click callFunctionOn: %w", err)
			}
			return nil
		}),
	)
}

// typeByNodeID focuses a DOM node and sends keyboard events.
func typeByNodeID(ctx context.Context, backendNodeID int64, text string) error {
	return chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			// scrollIntoView before focus
			objectID, err := resolveNodeToObject(ctx, backendNodeID)
			if err != nil {
				return err
			}
			scrollP := map[string]any{
				"objectId":            objectID,
				"functionDeclaration": "function() { this.scrollIntoViewIfNeeded(); }",
				"arguments":           []any{},
			}
			chromedp.FromContext(ctx).Target.Execute(ctx, "Runtime.callFunctionOn", scrollP, nil)

			p := map[string]any{"backendNodeId": backendNodeID}
			if err := chromedp.FromContext(ctx).Target.Execute(ctx, "DOM.focus", p, nil); err != nil {
				return fmt.Errorf("DOM.focus: %w", err)
			}
			return nil
		}),
		chromedp.KeyEvent(text),
	)
}

// resolveNodeToObject converts a backendNodeID to a JS remote object ID.
func resolveNodeToObject(ctx context.Context, backendNodeID int64) (string, error) {
	p := map[string]any{"backendNodeId": backendNodeID}
	var result json.RawMessage
	if err := chromedp.FromContext(ctx).Target.Execute(ctx, "DOM.resolveNode", p, &result); err != nil {
		return "", fmt.Errorf("DOM.resolveNode: %w", err)
	}
	var resp struct {
		Object struct {
			ObjectID string `json:"objectId"`
		} `json:"object"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("unmarshal resolveNode: %w", err)
	}
	if resp.Object.ObjectID == "" {
		return "", fmt.Errorf("no objectId for node %d", backendNodeID)
	}
	return resp.Object.ObjectID, nil
}
