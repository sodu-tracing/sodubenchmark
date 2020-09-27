package main

import (
	"fmt"
	"math/rand"
	v11 "sodubenchmark/go/common/v1"
	resoucev1 "sodubenchmark/go/resource/v1"
	v1 "sodubenchmark/go/trace/v1"
	"time"
)

func generateAttributes(n int) []*v11.KeyValue {
	kvs := make([]*v11.KeyValue, 0, n)
	for i := 0; i < n; i++ {
		kvs = append(kvs, &v11.KeyValue{
			Key: fmt.Sprintf("key%d", i),
			Value: &v11.AnyValue{
				Value: &v11.AnyValue_StringValue{
					StringValue: fmt.Sprintf("%d", rand.Int()),
				},
			},
		})
	}
	return kvs
}

func generateEvents(n int) []*v1.Span_Event {
	events := make([]*v1.Span_Event, 0, n)
	for i := 0; i < n; i++ {
		events = append(events, &v1.Span_Event{
			Name:         fmt.Sprintf("random logs at %d", n),
			TimeUnixNano: uint64(time.Now().UnixNano()),
			Attributes:   generateAttributes(3),
		})
	}
	return events
}

func generateTrace(depth int) []*v1.Span {
	trace := make([]*v1.Span, 0, depth)
	// generate traceid.
	traceID := make([]byte, 16)
	rand.Read(traceID)
	// temp variable for tracing previous trace id.
	prevSpanID := make([]byte, 16)
	prevSpanID = nil
	for depth > 0 {
		spanID := make([]byte, 16)
		rand.Read(spanID)
		startTime := time.Now()
		span := &v1.Span{
			TraceId:           traceID,
			SpanId:            spanID,
			ParentSpanId:      prevSpanID,
			Name:              fmt.Sprintf("function at %d", depth),
			Kind:              v1.Span_SPAN_KIND_SERVER,
			StartTimeUnixNano: uint64(startTime.UnixNano()),
			EndTimeUnixNano:   uint64(startTime.Add(time.Millisecond * 10).UnixNano()),
			Attributes:        generateAttributes(3),
			Events:            generateEvents(3),
		}
		trace = append(trace, span)
		prevSpanID = spanID
		depth--
	}
	return trace
}

func generateResourceAttributes(id int) []*v11.KeyValue {
	out := make([]*v11.KeyValue, 0, 2)
	out = append(out, &v11.KeyValue{
		Key: "instace",
		Value: &v11.AnyValue{
			Value: &v11.AnyValue_StringValue{StringValue: fmt.Sprintf("%d", id)},
		},
	})
	out = append(out, &v11.KeyValue{
		Key: "resource.name",
		Value: &v11.AnyValue{
			Value: &v11.AnyValue_StringValue{
				StringValue: fmt.Sprintf("%d", id),
			},
		},
	})
	return out
}

// generateSpans generates spans based on the number of spans that needs to
// generated. depth is used to specify how many spans per trace.
func generateTraces(num int, depth int) []*v1.ResourceSpans {
	resourceSpans := []*v1.ResourceSpans{}
	tmpSpans := make([][]*v1.Span, depth)
	for num > 0 {
		trace := generateTrace(depth)
		num -= len(trace)
		for i := 0; i < len(trace); i++ {
			tmpSpans[i] = append(tmpSpans[i], trace[i])
			// Batch all the resource spans if its' reaches 20
			if len(tmpSpans[i]) > 20 {
				resourceSpans = append(resourceSpans, &v1.ResourceSpans{
					Resource: &resoucev1.Resource{
						Attributes: generateResourceAttributes(i),
					},
					InstrumentationLibrarySpans: []*v1.InstrumentationLibrarySpans{
						{
							Spans: tmpSpans[i],
						},
					},
				})
				tmpSpans[i] = []*v1.Span{}
			}
		}
	}

	// Put remainging span.
	for i := 0; i < len(tmpSpans); i++ {
		if len(tmpSpans[i]) > 0 {
			resourceSpans = append(resourceSpans, &v1.ResourceSpans{
				Resource: &resoucev1.Resource{
					Attributes: generateResourceAttributes(i),
				},
				InstrumentationLibrarySpans: []*v1.InstrumentationLibrarySpans{
					{
						Spans: tmpSpans[i],
					},
				},
			})
			tmpSpans[i] = []*v1.Span{}
		}
	}
	return resourceSpans
}

func main() {

}
