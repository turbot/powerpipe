
benchmark "test_detection_benchmark_1" {
  title       = "Test Detection Benchmark"
  description = "This detection benchmark is used for testing."
  type        = "detection"
  children = [
    benchmark.test_detection_benchmark_2,
    benchmark.test_detection_benchmark_3,
  ]
}

benchmark "test_detection_benchmark_2" {
  title       = "Test Detection Benchmark 2"
  description = "This detection benchmark is used for testing."
  type        = "detection"
  children = [
    detection.test_detection_21,
    detection.test_detection_22,
  ]
}

benchmark "test_detection_benchmark_3" {
  title       = "Test Detection Benchmark 3"
  description = "This detection benchmark is used for testing."
  type        = "detection"
  children = [
    detection.test_detection_31,
    detection.test_detection_32,
  ]
}

detection "test_detection_21" {
  title           = "Test Detection 21"
  description     = "This detection is used for testing."
  severity        = "low"
  query           = query.detection_21
}

query "detection_21" {
  sql = <<-EOQ
    select 21 as detection_21;
  EOQ
}

detection "test_detection_22" {
  title           = "Test Detection 22"
  description     = "This detection is used for testing."
  severity        = "low"
  query           = query.detection_22
}

query "detection_22" {
  sql = <<-EOQ
    select 22 as detection_22;
  EOQ
}

detection "test_detection_31" {
  title           = "Test Detection 21"
  description     = "This detection is used for testing."
  severity        = "low"
  query           = query.detection_21
}

query "detection_31" {
  sql = <<-EOQ
    select 31 as detection_31;
  EOQ
}

detection "test_detection_32" {
  title           = "Test Detection 32"
  description     = "This detection is used for testing."
  severity        = "low"
  query           = query.detection_22
}

query "detection_32" {
  sql = <<-EOQ
    select 32 as detection_32;
  EOQ
}

benchmark "normal_benchmark" {
  title       = "Normal Benchmark"
  description = "This benchmark is used for testing."
  children = [
    control.normal1,
    control.normal2
  ]
}

control "normal1" {
  title       = "Normal Control 1"
  description = "This control is used for testing."
  type        = "control"
  severity    = "low"
  query       = query.normal1
}
query "normal1" {
  sql = <<-EOQ
    select 1 as normal1;
  EOQ
}
control "normal2" {
  title       = "Normal Control 2"
  description = "This control is used for testing."
  type        = "control"
  severity    = "low"
  query       = query.normal2
} 
query "normal2" {
  sql = <<-EOQ
    select 2 as normal2;
  EOQ
}