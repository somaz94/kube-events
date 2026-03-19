package event

import (
	"sort"
	"strings"
	"time"
)

// FilterOptions defines criteria for filtering events.
type FilterOptions struct {
	Since   time.Duration
	Kinds   []string
	Names   []string
	Types   []string
	Reasons []string
}

// Filter returns events matching all specified criteria.
func Filter(events []Event, opts FilterOptions) []Event {
	cutoff := time.Now().Add(-opts.Since)

	var result []Event
	for _, e := range events {
		if e.LastSeen.Before(cutoff) {
			continue
		}
		if len(opts.Kinds) > 0 && !containsCI(opts.Kinds, e.InvolvedObject.Kind) {
			continue
		}
		if len(opts.Names) > 0 && !containsCI(opts.Names, e.InvolvedObject.Name) {
			continue
		}
		if len(opts.Types) > 0 && !containsCI(opts.Types, e.Type) {
			continue
		}
		if len(opts.Reasons) > 0 && !containsCI(opts.Reasons, e.Reason) {
			continue
		}
		result = append(result, e)
	}

	// Sort by LastSeen descending (newest first)
	sort.Slice(result, func(i, j int) bool {
		return result[i].LastSeen.After(result[j].LastSeen)
	})

	return result
}

// GroupByResource groups events by their involved object.
func GroupByResource(events []Event) []ResourceGroup {
	m := make(map[ResourceKey][]Event)
	var keys []ResourceKey

	for _, e := range events {
		key := ResourceKey{
			Kind:      e.InvolvedObject.Kind,
			Name:      e.InvolvedObject.Name,
			Namespace: e.InvolvedObject.Namespace,
		}
		if _, exists := m[key]; !exists {
			keys = append(keys, key)
		}
		m[key] = append(m[key], e)
	}

	// Sort groups: resources with warnings first, then by newest event
	sort.Slice(keys, func(i, j int) bool {
		iWarn := hasWarning(m[keys[i]])
		jWarn := hasWarning(m[keys[j]])
		if iWarn != jWarn {
			return iWarn
		}
		return m[keys[i]][0].LastSeen.After(m[keys[j]][0].LastSeen)
	})

	groups := make([]ResourceGroup, len(keys))
	for i, key := range keys {
		groups[i] = ResourceGroup{Key: key, Events: m[key]}
	}
	return groups
}

// GroupEvents groups events by the specified mode.
func GroupEvents(events []Event, mode GroupBy) []ResourceGroup {
	if mode == GroupResource {
		return GroupByResource(events)
	}

	m := make(map[string][]Event)
	var keys []string

	for _, e := range events {
		var key string
		switch mode {
		case GroupNamespace:
			key = e.InvolvedObject.Namespace
			if key == "" {
				key = "(cluster-scoped)"
			}
		case GroupKind:
			key = e.InvolvedObject.Kind
		case GroupReason:
			key = e.Reason
		}
		if _, exists := m[key]; !exists {
			keys = append(keys, key)
		}
		m[key] = append(m[key], e)
	}

	// Sort: groups with warnings first, then by newest event
	sort.Slice(keys, func(i, j int) bool {
		iWarn := hasWarning(m[keys[i]])
		jWarn := hasWarning(m[keys[j]])
		if iWarn != jWarn {
			return iWarn
		}
		return m[keys[i]][0].LastSeen.After(m[keys[j]][0].LastSeen)
	})

	groups := make([]ResourceGroup, len(keys))
	for i, key := range keys {
		groups[i] = ResourceGroup{
			Key:    ResourceKey{Label: key},
			Events: m[key],
		}
	}
	return groups
}

func hasWarning(events []Event) bool {
	for _, e := range events {
		if e.Type == "Warning" {
			return true
		}
	}
	return false
}

func containsCI(list []string, target string) bool {
	for _, s := range list {
		if strings.EqualFold(s, target) {
			return true
		}
	}
	return false
}
